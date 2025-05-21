class CookieManager {
    static getCookie(name) {
        const cookies = document.cookie.split(";");
        for (const cookie of cookies) {
            const trimmed = cookie.trim();
            if (trimmed.startsWith(`${name}=`)) {
                return trimmed.substring(name.length + 1);
            }
        }
        return "";
    }

    static setCookie(name, value, days = 365) {
        const date = new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
        document.cookie = `${name}=${value};expires=${date.toUTCString()};path=/`;
    }
}
