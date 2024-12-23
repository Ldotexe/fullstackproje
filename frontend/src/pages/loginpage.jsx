import {Flex, Typography} from "antd";
import Login from "../components/auth/login.jsx";
import SSLogo from "../assets/SS.png";

const {Title} = Typography;

const LoginPage = () => {
    return (
        <Flex vertical align={"center"} gap={"large"} style={{height: "100vw"}}>
            <img src={SSLogo} alt="SS" style={{width: "auto", height: "40vh"}}/>
            <Login/>
        </Flex>
    );
}

export default LoginPage;