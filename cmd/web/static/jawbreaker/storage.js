class LocalStorageManager {
    static getItem(name) {
        const value = localStorage.getItem(name);
        return value || "";
    }

    static setItem(name, value) {
        localStorage.setItem(name, value);
    }
}