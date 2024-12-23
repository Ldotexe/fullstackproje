import "./App.css";
import {Layout, Typography, message} from "antd";
import {createContext, useEffect, useState} from "react";
import {Route, Routes} from "react-router-dom";
import LoginPage from "./loginpage.jsx";
import MainPage from "./mainpage.jsx";
import ChatPage from "./chatpage.jsx";

const {Text, Title} = Typography;
const {Content, Header} = Layout;

export const MessageContext = createContext(null);

const contentStyle = {
    textAlign: 'center',
    // minHeight: 120,
    lineHeight: '120px',
    // color: '#fff',
    // backgroundColor: '#efefef00',
};

const App = () => {
    const [messageApi, contextHolder] = message.useMessage();

    const show_error = (text) => {
        messageApi.error(text);
    };

    const show_success = (text) => {
        messageApi.success(text);
    };

    return (
        <MessageContext.Provider value={[show_error, show_success]}>
            <Layout style={{
                overflowY: 'hidden',
                overflowX: 'hidden',
                width: 'calc(100%)',
                maxWidth: 'calc(100%)',
                height: '100vh',
            }}>
                <Content style={contentStyle}>
                    <Routes>
                        <Route exact path={"/login"} element={<LoginPage/>}/>
                        <Route exact path={"/"} element={<MainPage/>}/>
                        <Route exact path={"/chat/:username"} element={<ChatPage/>}/>
                    </Routes>
                </Content>
                {contextHolder}
            </Layout>
        </MessageContext.Provider>
    );
}

export default App;
