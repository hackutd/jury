import { useEffect, useState } from 'react';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import { useNavigate } from 'react-router-dom';

const Judge = () => {
    const navigate = useNavigate();
    const [verified, setVerified] = useState(false);
    const [prev, setPrev] = useState(null);
    const [next, setNext] = useState(null);

    // Verify user is logged in and read welcome before proceeding
    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedIn = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/auth`, {
                method: 'POST',
                credentials: 'include',
            });
            if (!loggedIn.ok) {
                console.error(`Judge is not logged in! ${loggedIn.status} ${loggedIn.statusText}`);
                navigate('/judge/login');
                return;
            }

            // Check for read welcome
            const readWelcomeRes = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/welcome`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            });
            if (!readWelcomeRes.ok) {
                alert(
                    `Unable to connect to server: ${readWelcomeRes.status} ${readWelcomeRes.statusText}. Please check your connection or reload the page.`
                );
                return;
            }
            const readWelcome = await readWelcomeRes.json();
            if (!readWelcome.ok) {
                navigate('/judge/welcome');
            }

            setVerified(true);
        }

        fetchData();
    }, []);

    useEffect(() => {
        if (!verified) return;

        // Once verification finishes, get the judge's next and prev project to judge
    }, [verified]);

    return (
        <>
            <JuryHeader withLogout />
            <Container>Judging Page!</Container>
        </>
    );
};

export default Judge;
