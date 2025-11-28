import { Link } from "react-router-dom"

export const SuccessView = () => {
    return (
        <div>
            <h1>Заказ успешно оформлен</h1>
            <Link to="/" className="nav-link"><button>На главную</button></Link>
        </div>
    )
}