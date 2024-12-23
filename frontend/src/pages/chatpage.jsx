import {useNavigate, useParams} from "react-router-dom";
import UserChat from "../components/chat/userchat.jsx";
import {postSendMessage, useMessages} from "../hooks/api.jsx";
import {useContext, useEffect, useState} from "react";
import {MessageContext} from "./App.jsx";


const ChatPage = () => {
    const {username} = useParams();
    const [requestData, setRequestData] = useState(null);
    const [messages, isMessagesLoading, messagesError] = useMessages(username, [requestData]);
    const [showError, showSuccess] = useContext(MessageContext);
    const navigate = useNavigate();

    useEffect(() => {
        if (messagesError) {
            showError("No auth");
            navigate("/");
        }
    }, [messagesError]);

    useEffect(() => {
        if (requestData) {
            setRequestData(null);
        }
    }, [requestData]);


    const sendMessage = (username, text) => {
        postSendMessage(username, text, setRequestData);
    }
    

    return (
        <UserChat username={username} messages={messages ? messages : []} sendMessage={sendMessage}/>

    );
}

export default ChatPage;