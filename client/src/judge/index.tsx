import { useEffect } from 'react';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';

const Judge = () => {
    useEffect(() => {
        // TODO: Check to see if a user is already logged in
        // TODO: Check to see if a user has already read welcome page
    }, []);

    return (
        <>
            <JuryHeader withLogout />
            <Container>Judging Page!</Container>
        </>
    );
};

export default Judge;
