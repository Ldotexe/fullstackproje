import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import {BrowserRouter} from "react-router-dom";
import App from './pages/App.jsx'
import {ConfigProvider, theme} from "antd";

const main_color = "#48cf68";

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <BrowserRouter>
            <ConfigProvider
                theme={{
                    // 1. Use dark algorithm
                    algorithm: theme.darkAlgorithm,
                    components: {
                        Button: {
                            colorLink: main_color,
                            colorPrimary: main_color,
                            algorithm: true, // Enable algorithm
                        },
                        Input: {
                            colorLink: main_color,
                            colorPrimary: main_color,
                            algorithm: true, // Enable algorithm
                        },
                        Typography: {

                            colorLink: main_color,
                            algorithm: true, // Enable algorithm
                        },
                        Select: {
                            colorLink: main_color,
                            colorPrimary: main_color,
                            algorithm: true
                        }
                    },


                    // 2. Combine dark algorithm and compact algorithm
                    // algorithm: [theme.darkAlgorithm, theme.compactAlgorithm],
                }}
            >
                <App/>
            </ConfigProvider>
        </BrowserRouter>
    </StrictMode>,
)
