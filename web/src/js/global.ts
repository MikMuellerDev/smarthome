// Global data interface
export interface Data {
    userData: UserData
    notifications: Notification[]
    notificationCount: number
    notificationsLoaded: boolean
    notificationDoneMarkerAdded: boolean
}

// User data fetched from the server
export interface UserData {
    username: string
    forename: string
    surname: string
    primaryColor: string
}

export interface Notification {
    id: number
    priority: number
    name: string
    description: string
    date: string
}

// Global datastore, it populated when the page loads
export var data: Data = {
    userData: {
        forename: "",
        primaryColor: "",
        surname: "",
        username: "",
    },
    notifications: [],
    notificationCount: 0,
    notificationsLoaded: false,
    notificationDoneMarkerAdded: false
}

export async function fetchData() {
    data.userData = await fetchUserData()
    data.notificationCount = await fetchNotificationCount()
}

export async function fetchUserData(): Promise<UserData> {
    return await (await fetch('/api/user/data')).json()
}

export async function fetchNotifications(): Promise<Notification[]> {
    return await (await fetch('/api/user/notification/list')).json()
}

export async function fetchNotificationCount(): Promise<number> {
    return await (await fetch('/api/user/notification/count')).json()
}

export const sleep = (ms: number) => new Promise((res) => setTimeout(res, ms))
