import {useEffect, useState} from "react";
import axios from "axios";
import {apiURL} from "../config";


export const useApi = (url, request_data = {}, deps = []) => {
    const [data, setData] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);


    useEffect(() => {
        setIsLoading(true);

        axios
            .get(apiURL + url, {params: request_data})
            .then((res) => {
                console.log(res.data);
                setData(res.data);
                setIsLoading(false);
            })
            .catch(err => {
                setError(err.message);
                setIsLoading(false);
            });
    }, deps);

    return [data, isLoading, error];
}

export const useMessages = (username, deps = []) => useApi("/messages/" + username, {}, deps);
export const useUsers = (deps = []) => useApi("/users", {}, deps);


export const postApi = (url, data, setData = ((e) => {
                        }), setError = ((e) => {
                        }),
                        reload = false) => {
    axios
        .post(apiURL + url, data)
        .then((res) => {
            console.log(res.data);
            setData(res.status);

            return res.data;
        })
        .catch(err => {
            console.log(err.message);
            setError(err.message);

            return err.message;
        });

}

export const postAuth = (login, password, setData) => postApi(
    "/auth",
    {login, password},
    setData,
    setData
);

export const postSendMessage = (username, text, setData) => postApi(
    "/message/send/" + username,
    {text},
    setData,
    setData
);