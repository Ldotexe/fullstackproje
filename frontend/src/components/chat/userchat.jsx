import {Avatar, Divider, Flex, Typography, Input, Button, Empty} from "antd";
import {CaretLeftOutlined, CaretRightOutlined, RollbackOutlined, SendOutlined} from "@ant-design/icons";
import {useEffect, useRef, useState} from "react";
import {postSendMessage} from "../../hooks/api.jsx";
import {useNavigate} from "react-router-dom";

const {Text} = Typography;
const {TextArea} = Input;

const ChatMessage = ({text, isMine}) => {
    return (
        <Flex justify={"flex-start"} align={"center"} gap={"large"} style={{width: "90vw"}}>
            {isMine ? <Avatar size={"large"} style={{backgroundColor: "#6acf48"}}><CaretRightOutlined/></Avatar> :
                <Avatar size={"large"} style={{backgroundColor: "#4897cf"}}><CaretLeftOutlined/></Avatar>}
            <Text style={{width: "90%", textAlign: "left"}}>{text}</Text>
        </Flex>
    );
}

const ChatInput = ({username, sendMessage}) => {
    const [messageText, setMessageText] = useState("");

    const onSend = () => {
        console.log("sending message to " + username + ": " + messageText);
        sendMessage(username, messageText);
    }

    return (
        <Flex gap={"middle"} justify={"flex-start"} align={"center"}
              style={{marginTop: "5%", bottom: "3vh", width: "90%"}}>
            <TextArea
                onChange={(e) => setMessageText(e.target.value)}
                style={{
                    // maxWidth: "100vh",
                    width: "100%"
                }}
                onPressEnter={onSend}
                placeholder="Введите ваше сообщение.."
                showCount
                maxLength={1500}
                autoSize={{
                    minRows: 3,
                    maxRows: 3,
                }}

            />
            <Button shape={"circle"} icon={<SendOutlined/>} size={"large"}
                    onClick={onSend}></Button>
        </Flex>

    );
}


const UserChat = ({username, messages, sendMessage}) => {
    const navigate = useNavigate();
    const messagesEndRef = useRef(null);
    console.log(messages);

    const scrollToBottom = () => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollTop = messagesEndRef.current.scrollHeight;
        }
    };

    useEffect(() => {
        console.log("Scrolling to bottom");
        scrollToBottom();
    }, [messages]);


    return (
        <Flex vertical gap={"large"} align={"center"} justify={"space-evenly"}
              style={{marginTop: 0, overflowY: 'hidden', overflowX: 'hidden', height: '100%', width: "100%"}}>
            <Flex gap={"middle"} align={"center"}>
                <Avatar size={"large"}>{username.slice(0, 1).toUpperCase()}</Avatar>
                <Text style={{fontSize: "20pt"}}>{username}</Text>
                <Button size={"large"} type={"link"} icon={<RollbackOutlined />} onClick={() => navigate("/")}></Button>
            </Flex>
            <Flex ref={messagesEndRef} vertical gap={"small"} style={{
                overflowY: 'scroll',
                overflowX: 'hidden',
                height: '65vh',
                width: "90%",
                paddingTop: "7%"
            }}>
                <Divider/>
                {(messages && messages.length === 0) && <Empty description={"Нет сообщений"}/>}
                {messages.map((message) => <><ChatMessage text={message.text} isMine={message.isMine}/><Divider/></>)}
            </Flex>
            <ChatInput username={username} sendMessage={sendMessage}/>
        </Flex>
    );
}

export default UserChat;