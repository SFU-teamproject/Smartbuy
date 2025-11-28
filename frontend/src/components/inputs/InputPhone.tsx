import { ChangeEvent, useRef } from 'react';
import { IMaskInput } from 'react-imask';

interface InputProps {
    number: string;
    setNumber: React.Dispatch<React.SetStateAction<string>>;
    setIsValid: React.Dispatch<React.SetStateAction<boolean>>;
    setErrorPhoneNumber: React.Dispatch<React.SetStateAction<string>>;
};

const InputPhone: React.FC<InputProps> = ({number, setNumber, setIsValid, setErrorPhoneNumber}: InputProps) => {
    const ref = useRef(null);
    const inputRef = useRef(null);

    const handlePhone = (event: ChangeEvent<HTMLInputElement>) => {
            setNumber(event.target.value);
        }

    const validatePhone = () => {
        if (!number) {
            setErrorPhoneNumber("Number must be not empty");
            setIsValid(false);
            return false;
        }
        //const regexp = /^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$/;
        const regexp = /^\+?7[\s]?[(]?\d{3}[)]?[\s]?\d{3}-?\d{2}-?\d{2}$/;
        if (!regexp.test(number)) {
            setErrorPhoneNumber("Not accepted");
            setIsValid(false);
            return false;
        }
        setErrorPhoneNumber("");
        setIsValid(true);
        return true;
    }

    return (
        <div>
            <IMaskInput
            mask="+7 (000) 000-00-00"
            unmask={true} // true|false|'typed'
            ref={ref}
            inputRef={inputRef}  // access to nested input
            // DO NOT USE onChange TO HANDLE CHANGES!
            // USE onAccept INSTEAD
            //onAccept={
                // depending on prop above first argument is
                // `value` if `unmask=false`,
                // `unmaskedValue` if `unmask=true`,
                // `typedValue` if `unmask='typed'`
              //  (value, mask) => console.log(value)
            //}
            onBlur={validatePhone}
            onChange={handlePhone}
            // ...and more mask props in a guide
            // input props also available
            placeholder='+7 (___) ___-__-__'
            className='input-text'
            >
            </IMaskInput>
        </div>
    );
} 
export default InputPhone;
