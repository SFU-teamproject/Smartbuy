import './Payment.css';
import { useAuth } from '../../context/AuthContext';
import { getCartItems, getSmartphoneById } from '../../api/client';
import { CartItem, Smartphone, Payment, DeliveryType, PaymentType } from '../../types';
import { ChangeEvent, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import InputPhone from '../inputs/InputPhone';

export const PaymentView = () => {
    const { cart, token } = useAuth();
    const [items, setItems] = useState<(CartItem & { smartphone?: Smartphone })[]>([]);
    const [loading, setLoading] = useState(true);
    const [total, setTotal] = useState(0);
    const [paymentType, setPaymentType] = useState<PaymentType>(PaymentType.payment_self);
    const [deliveryType, setDeliveryType] = useState <DeliveryType>(DeliveryType.delivery_courier);
    const [city, setCity] = useState("");
    const [street, setStreet] = useState("");
    const [house, setHouse] = useState("");
    const [errorPhoneNumber, setErrorPhoneNumber] = useState("");
    const [errorCity, setErrorCity] = useState("");
    const [errorStreet, setErrorStreet] = useState("");
    const [errorHouse, setErrorHouse] = useState("");
    const [phoneNumber, setPhoneNumber] = useState("");
    const [phoneIsValid, setPhoneIsValid] = useState(false);
    const navigate = useNavigate();
    const handleSubmit = () => {
        if (deliveryType === DeliveryType.delivery_courier && !validatePaymentForm() ) {
            return
        }
        if (deliveryType === DeliveryType.delivery_self && !phoneIsValid ) {
            setErrorPhoneNumber("Number must not be null");
            return
        }
        /* для отправки */
        const paymentData: Payment = {
            phone: phoneNumber,
            delivery_type: deliveryType,
            address: {
                city: city,
                street: street,
                house: parseInt(house)
            },
            payment_type: paymentType
        };
        console.log(paymentData);
        console.log(errorPhoneNumber);
        (paymentType === PaymentType.payment_cart) ? navigate("/bank") : navigate("/success")
    }
    
    useEffect(() => {
        if (!cart?.id || !token) {
            setLoading(false);
            return;
        }

        const loadCartData = async () => {
            try {
                const cartItems = await getCartItems(cart.id, token);

                // Загружаем информацию о каждом товаре
                const itemsWithProducts = await Promise.all(
                    cartItems.map(async (item) => {
                        try {
                            const smartphone = await getSmartphoneById(item.smartphone_id, token);
                            return { ...item, smartphone };
                        } catch {
                            return item; // Если не удалось загрузить данные товара
                        }
                    })
                );

                setItems(itemsWithProducts);
            } catch (error) {
                console.error('Ошибка загрузки корзины:', error);
            } finally {
                setLoading(false);
            }
        };

        loadCartData();
    }, [cart?.items, token]);

    // Подсчет общей суммы
    useEffect(() => {
        const newTotal = items.reduce((sum, item) => {
            const price = item.smartphone?.price || 0;
            return sum + (price * item.quantity);
        }, 0);
        setTotal(newTotal);
    }, [items]);

    const handleStreet  = (event: ChangeEvent<HTMLInputElement>) => {
        setStreet(event.target.value);
    }
    const handleHouse = (event: ChangeEvent<HTMLInputElement>) => {
        setHouse(event.target.value);
    }
    const handleCity = (event: ChangeEvent<HTMLInputElement>) => {
        setCity(event.target.value);
    }
    const handleDeliveryType = (event: ChangeEvent<HTMLInputElement>) => {
        const typeValue = event.target.value;
        if((Object.values(DeliveryType) as string[]).includes(typeValue)) 
            setDeliveryType(typeValue as unknown as DeliveryType)
    }
    const handlePaymentType = (event: ChangeEvent<HTMLInputElement>) => {
        const typeValue = event.target.value;
        if((Object.values(PaymentType) as string[]).includes(typeValue)) 
            setPaymentType(typeValue as unknown as PaymentType)
    }
    const validateCity = () => {
        if (!city || city.length <= 2) {
            setErrorCity("City length must be > 2 and not empty");
            return false;
        }
        const regexp = /^[A-Za-zА-Яа-яё]*$/;
        if (!regexp.test(city)) {
            setErrorCity("Only letters");
            return false;
        }
        setErrorCity("");
        return true;
    }
    const validateStreet = () => {
        if (!street || street.length <= 2) {
            setErrorStreet("Street length must be > 2 and not empty");
            return false;
        }
        const regexp = /^[A-Za-zА-Яа-яё]*$/;
         if (!regexp.test(street)) {
            setErrorStreet("Only letters");
            return false;
         }
        setErrorStreet("");
        return true;
    }
    const validateHouse = () => {
        if (!house) {
            setErrorHouse("House input must be not empty");
            return false;
        }
        const regexp = /^\d{1,3}$/;
        if (!regexp.test(house)) {
            setErrorHouse("Only numbers. length must be >= 1 & <= 3");
            return false;
        }
        setErrorHouse("");
        return true;
    }

    const validatePaymentForm = (): boolean => {
        const vCity = validateCity();
        const vStreet = validateStreet();
        const vHouse = validateHouse();
        if (!phoneIsValid && !errorPhoneNumber)
            setErrorPhoneNumber("Number must not be null");
        return vCity && vStreet && vHouse && phoneIsValid;
    }

    if (loading) return <div className="loading">Загрузка...</div>;

    return (
        <div className="payment_container">
            <h1>Ваш заказ</h1>
            <div className='payment_bucket'>
                {items.sort((a, b) => a.id - b.id).map(item => (
                    <div className="product-info">
                        <span className='product-info-title'>{item.smartphone?.producer} {item.smartphone?.model || `Товар #${item.smartphone_id}`}</span>
                        <div className='product-info-numbers'>
                            <span className='product-info-quantity'>{item.quantity} шт.</span>
                            <span className="product-info-price">{item.smartphone?.price ? `${(item.smartphone.price * item.quantity).toLocaleString('ru-RU')} ₽` : '—'}</span>
                        </div>
                    </div>
                ))}
                <div className='payment_right-total'>
                    <strong>Итого:</strong>
                    <strong>{total.toLocaleString()} ₽</strong>
                </div>
            </div>
            <h1>Оформление заказа</h1>
            <div className="payment_order">
                <div className='payment_order-box'>
                    <h2>Ваш номер телефона</h2>
                    <label htmlFor="phone"></label>
                    <InputPhone setNumber={setPhoneNumber} number={phoneNumber} setIsValid={setPhoneIsValid} setErrorPhoneNumber={setErrorPhoneNumber}></InputPhone>
                    {errorPhoneNumber && (<div className='payment-error'>{errorPhoneNumber}</div>)}
                    <br />
                </div>
                <div className='payment_order-box'>
                    <h2>Доставка</h2>
                    <div className="form_radio_group">
                        <div className="form_radio_group-item">
                            <input type='radio' name='courier_type' value={DeliveryType.delivery_self} id='delivery_self' checked={deliveryType === DeliveryType.delivery_self} onChange={handleDeliveryType}></input>
                            <label htmlFor='delivery_self'>Самовывоз</label>
                        </div>
                        <div className="form_radio_group-item">
                            <input type='radio' name='courier_type' value={DeliveryType.delivery_courier} id='delivery_courier' checked={deliveryType === DeliveryType.delivery_courier}  onChange={handleDeliveryType}></input>
                            <label htmlFor='delivery_courier'>Курьером</label>
                        </div>
                    </div>
                    { deliveryType === DeliveryType.delivery_courier &&  (
                        <div className='payment_order-city'>
                            <div className='city-container'>
                                <label htmlFor="city">Город</label>
                                <input id="city" className='input-text' onBlur={validateCity} onChange={handleCity}></input>
                                {errorCity && (<div className='payment-error'>{errorCity}</div>)}
                            </div>
                            <div className='city-container'>
                                <label htmlFor="street">Улица</label>
                                <input id="street" className='input-text' onBlur={validateStreet} onChange={handleStreet}></input>
                                {errorStreet && (<div className='payment-error'>{errorStreet}</div>)}
                            </div>
                            <div className='city-container'>
                                <label htmlFor="house_number">Номер дома</label>
                                <input id="house_number" className='input-text' onBlur={validateHouse} onChange={handleHouse}></input>
                                {errorHouse && (<div className='payment-error'>{errorHouse}</div>)}
                            </div>
                        </div>
                    )} 
                </div>
                { deliveryType === DeliveryType.delivery_courier && (
                    <div className='payment_order-box'>
                        <h2>Способ оплаты</h2>
                        <div className="form_radio_group">
                            <div className="form_radio_group-item">
                                <input type="radio" name="payment_type" id="payment_cart" value={PaymentType.payment_cart} checked={paymentType === PaymentType.payment_cart} onChange={handlePaymentType}></input>
                                <label htmlFor="payment_cart">Картой</label>
                            </div>
                            <div className="form_radio_group-item">
                                <input type="radio" name="payment_type" id="payment_courier" value={PaymentType.payment_courier} checked={paymentType === PaymentType.payment_courier} onChange={handlePaymentType}></input>
                                <label htmlFor="payment_courier">Курьер</label>
                            </div>
                            <div className="form_radio_group-item">
                                <input type="radio" name="payment_type" id="payment_self" value={PaymentType.payment_self} checked={paymentType === PaymentType.payment_self} onChange={handlePaymentType}></input>
                                <label htmlFor="payment_self">Наличными</label>
                            </div>
                        </div>
                    </div>
                )}
            </div>
            <button className='submit' type="submit" onClick={handleSubmit}>Оформить заказ</button>
        </div>
    );
}