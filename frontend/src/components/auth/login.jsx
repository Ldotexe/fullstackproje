import {Flex, Form, Input, Typography, Button} from "antd";
import {postAuth} from "../../hooks/api.jsx";
import {useContext, useEffect, useState} from "react";
import {MessageContext} from "../../pages/App.jsx";
import {useNavigate} from "react-router-dom";

const {Title} = Typography;

const Login = () => {
    const navigate = useNavigate();
    const [requestData, setRequestData] = useState(null);
    const [showError, showSuccess] = useContext(MessageContext);

    useEffect(() => {
        if (requestData && requestData === 200) {
            console.log("uraaa");
            showSuccess("SUCCESS!!!");
            navigate("/");
        }
    }, [requestData]);

    return (
        <Flex vertical justify={"center"} align={"center"}>
            <Title level={5}>Авторизация/Регистрация</Title>
            <Form
                name="auth"
                style={{
                    minWidth: "80vw",
                }}
                layout={"vertical"}
                onFinish={(e) => {
                    console.log("success", e);
                    postAuth(e.login, e.password, setRequestData);
                }}
                onFinishFailed={(e) => {
                    console.log("fail", e);
                    showError("FAIL!");
                }}
                autoComplete="off"
            >
                <Form.Item
                    label="Логин"
                    name="login"
                    rules={[
                        {
                            required: true,
                            message: 'Введите логин',
                        },
                        {
                            pattern: "^[a-zA-Z0-9_]*$",
                            message: 'Разрешены только символы a-z, A-Z, 0-9'
                        }
                    ]}
                >
                    <Input placeholder="johnappleseed"
                           style={{width: "80vw", maxWidth: 400}}
                           allowClear/>
                </Form.Item>
                <Form.Item
                    label="Пароль"
                    name="password"
                    rules={[
                        {
                            required: true,
                            message: 'Введите пароль',
                        },
                        {
                            pattern: "^(?=.*\\d).{8,}$",
                            message: 'Минимум 8 символов и минимум 1 цифра'
                        }
                    ]}
                >
                    <Input.Password style={{width: "80vw", maxWidth: 400}}/>
                </Form.Item>
                <Button size={"large"} autoFocus htmlType={"submit"} type={"primary"}
                        style={{width: "80vw", maxWidth: 400}}>Войти</Button>

            </Form>
        </Flex>
    );
}

export default Login;