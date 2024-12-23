import {Avatar, Flex, Card, Typography, Empty} from "antd";
import {useNavigate} from "react-router-dom";

const {Text} = Typography;

const Chat = ({name}) => {
    const navigate = useNavigate();

    return (
        <Card>
            <a onClick={() => navigate("/chat/" + name)}>
                <Flex justify={"flex-start"} align={"center"} gap={"middle"}>
                    <Avatar>{name.slice(0, 1).toUpperCase()}</Avatar>
                    <Text>{name}</Text>
                </Flex>
            </a>
        </Card>
    );
}

const Chats = ({chats}) => {
    if (chats && chats.length === 0) {
        return (
            <Empty description={"У вас нет чатов"}/>
        );
    }
    return (
        <Flex vertical gap={"small"} style={{width: "80vw"}}>
            {chats.map((username) => <Chat key={username} name={username}/>)}
        </Flex>
    );
}

export default Chats;