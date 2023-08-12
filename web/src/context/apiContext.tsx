import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { Api } from '@concurrent-world/client'

interface ApiContextState {
    api: Api
    setJWT: (jwt: string) => void
}

const defaultApiContextState: ApiContextState = {
    api: new Api({host: ''}),
    setJWT: (_) => { }
}

const ApiContext = createContext<ApiContextState>(defaultApiContextState)

export interface ApiProviderProps {
    children: JSX.Element
}

export default function ApiProvider(props: ApiProviderProps): JSX.Element {

    const [api, setApi] = useState<Api>(new Api({host: ''}))
    const [token, setToken] = useState<string | undefined>(undefined)

    useEffect(() => {
        if (!token) return
        setApi(new Api({host: '', token: token}))
    }, [])

    const setJWT = useMemo(() => (jwt: string) => {
        setToken(jwt)
        setApi(new Api({host: '', token: jwt}))
    }, [setApi])

    const apiContextState = useMemo(() => ({
        api,
        setJWT
    }), [api, setJWT])

    return <ApiContext.Provider value={apiContextState}>{props.children}</ApiContext.Provider>
}

export function useApi(): ApiContextState {
    return useContext(ApiContext)
}