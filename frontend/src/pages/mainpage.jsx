import {Flex, Typography} from "antd";
import Chats from "../components/chats/chats.jsx";
import {useUsers} from "../hooks/api.jsx";
import {useEffect} from "react";
import {useNavigate} from "react-router-dom";

const {Title} = Typography;

const MainPage = () => {
    const navigate = useNavigate();
    const [users, isUsersLoading, usersError] = useUsers();

    useEffect(() => {
        if (usersError) {
            navigate("/login");
        }
    }, [users, usersError]);

    return (
        <Flex vertical align={"center"}>
            <Title level={4}>Чаты</Title>
            <Chats chats={users ? users : []}/>
        </Flex>
    );
}

export default MainPage;