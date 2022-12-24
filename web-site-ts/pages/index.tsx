/* eslint-disable jsx-a11y/alt-text */
/* eslint-disable @next/next/no-img-element */
import { signIn, signOut, useSession } from "next-auth/react";

const Index = () => {
    const { data: session } = useSession();

    if (session) {
        return (
            <>
                Signed in as <img src={session.user?.image ?? ""} width="50px" /> <br />
                Name : {session.user?.name ?? ""} <br />
                AccessToken : {session.accessToken} <br />
                <button onClick={() => signOut()}>Sign out</button>
            </>
        );
    }

    return (
        <>
            Not signed in <br />
            <button onClick={() => signIn()}>Sign in</button>
        </>
    );
};

export default Index;
