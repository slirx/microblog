// todo move to options
const baseURL = "http://localhost:8080";

export default {
    install(Vue, options) {
        Vue.prototype.$http = {
            get: (url, accessToken, callback, isLoaderUsed) => {
                if (isLoaderUsed) {
                    options.store.commit('isLoading', true);
                }

                fetch(baseURL + url, {
                    method: "GET",
                    headers: {
                        'Authorization': accessToken
                    }
                }).then((response) => {
                    if (response.status === 403) {
                        localStorage.removeItem("jwt")
                        location.href = "/"; // todo do it without redirect
                        return;
                    }

                    if (response.status === 404) {
                        if (isLoaderUsed) {
                            options.store.commit('isLoading', false);
                        }

                        return;
                    }

                    if (!response.ok) {
                        console.log("!response.ok");
                        // TODO show specific error message with "reconnect"/"try again in n seconds" buttons
                        throw Error(response.statusText);
                    }

                    response.json().then((data) => {
                        callback(data);
                    });

                    if (isLoaderUsed) {
                        options.store.commit('isLoading', false);
                    }
                }).catch((error) => {
                    console.log(error);
                    options.store.commit('alert', {
                        type: "network-error",
                        message: "something went wrong. please, try later or contact support@microblog.local"
                    });
                    //options.store.commit('isLoading', false);
                });
                // .finally(() => {
                //         options.store.commit('isLoading', false);
                //     })

            },
            sendRequest: (url, method, accessToken, formData, callbackSuccess, callbackError, callbackFinally, isLoaderUsed) => {
                if (isLoaderUsed) {
                    options.store.commit('isLoading', true);
                }

                fetch(baseURL + url, {
                    method: method,
                    headers: {
                        'Authorization': accessToken
                    },
                    body: formData
                }).then((response) => {
                    if (!response.ok) {
                        // todo display also request id here
                        response.json().then((data) => {
                            options.store.commit('alert', {
                                type: data.type,
                                message: data.message
                            });
                        });

                        if (typeof callbackError != 'undefined') {
                            callbackError(response);
                        }

                        // TODO show specific error message with "reconnect"/"try again in n seconds" buttons
                        //throw Error(response.statusText);
                        return;
                    }

                    // todo write data to alert for errors

                    response.json().then((data) => {
                        // if user is not logged in
                        if (typeof data.code !== "undefined" && data.code === "not_logged_in") {
                            location.href = "/";
                            return;
                        }

                        if (typeof callbackSuccess !== 'undefined') {
                            callbackSuccess(data);
                        }
                    });
                }).catch((error, response) => {
                    console.log("error happened:", error);
                }).finally(() => {
                    if (typeof callbackFinally != 'undefined') {
                        callbackFinally();
                    }
                    if (isLoaderUsed) {
                        options.store.commit('isLoading', false);
                    }
                });
            }
        }
    }
};
