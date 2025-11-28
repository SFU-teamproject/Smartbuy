import { Link } from "react-router-dom"
import './Bank.css';

export const BankView = () => {
    return (
        <div className="bank_container">
            <h1 className="bank_title">Карта</h1>
            <div className="form-container">
                <div className="field-container">
                    <label htmlFor="cardholder">Name</label>
                    <input id="cardholder" maxLength={20} type="text"/>
                </div>
                <div className="field-container">
                    <label htmlFor="cardnumber">Card Number</label>
                    <input id="cardnumber" type="text" maxLength={16} pattern="[0-9]*" inputMode="numeric"/>
                </div>
                <div className="field-container">
                    <label htmlFor="cardexpirationdate">Expiration (mm/yy)</label>
                    <input id="cardexpirationdate" type="text" maxLength={4} pattern="[0-9]*" inputMode="numeric"/>
                </div>
                <div className="field-container">
                    <label htmlFor="cardsecuritycode">Security Code</label>
                    <input id="cardsecuritycode" type="text" maxLength={3} pattern="[0-9]*" inputMode="numeric"/>
                </div>
            </div>

            <Link to="/success" className="nav-link"><button>Оплатить заказ</button></Link>
        </div>
    )
}