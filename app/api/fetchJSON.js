var API = (function(API){

    API.ROUTE = "/api";

    API.fetchJSON = (route, callback) => {
        return fetch(API.ROUTE+route)
            .then((r) => r.json())
            .then((data) => callback(data));
    };

    return API;
}(API || {}));
