/* eslint-disable jsx-a11y/alt-text */
/* eslint-disable @next/next/no-img-element */
import { signOut, useSession } from "next-auth/react";
import { useRouter } from "next/router";

const Vtsm = () => {
    const { data: session } = useSession();
    const router = useRouter();
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

    router.push("/");
};

export default Vtsm;
