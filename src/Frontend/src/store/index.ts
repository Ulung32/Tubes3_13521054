import { Chat, History, QnA } from "@/models"
import { create } from "zustand"

type AuthState = {
    id: string
    username: string

    actions: {
        setId: (id: string) => void
        setUsername: (username: string) => void
    }
}

export const useAuthStore = create<AuthState>((set) => ({
    id: "",
    username: "",

    actions: {
        setId(id){
            set({id})
        },
        setUsername(username){
            set({username})
        }
    }
}))


export function useAuthAction(){
    return useAuthStore((state) => state.actions)
}


type ChatState = {
    idHistorySelected: string
    chats: Chat[]
    histories: History[]
    qnas: QnA[]
    algo: string
    isOpenSidebar: boolean

    actions: {
        setIdHistorySelected: (idHistorySelected: string) => void
        setChats: (chats: Chat[]) => void
        setHistories: (histories: History[]) => void
        setAlgo: (algo: string) => void
        setQnas: (qnas: QnA[]) => void
        setIsOpenSidebar: (isOpenSidebar: boolean) => void
    }
}

export const useChatStore = create<ChatState>((set) => ({
    idHistorySelected: "",
    chats: [],
    histories: [],
    qnas: [],
    algo: "kmp",        
    isOpenSidebar: false,

    actions: {
        setIdHistorySelected(idHistorySelected){
            set({idHistorySelected})
        },
        setChats(chats){
            set({chats})
        },
        setHistories(histories){
            set({histories})
        },
        setAlgo(algo){
            set({algo})
        },
        setQnas(qnas) {
            set({qnas})
        },
        setIsOpenSidebar(isOpenSidebar) {
            set({isOpenSidebar})
        },
    }
}))


export function useChatAction(){
    return useChatStore((state) => state.actions)
}

