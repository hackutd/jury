import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import Cookies from 'universal-cookie';
import Button from '../../components/Button';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import PasswordInput from '../../components/PasswordInput';
import Loading from '../../components/Loading';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';

const JudgeLogin = () => {
    const [searchParams] = useSearchParams();
    const [codeParam] = useState(searchParams.get('code') || '');

    const [disabled, setDisabled] = useState(true);
    const [code, setCode] = useState('');
    const [loginLock, setLoginLock] = useState(false);
    const [error, setError] = useState(false);
    const navigate = useNavigate();

    // If token cookie is already defined and valid, redirect to judge page
    useEffect(() => {
        if (codeParam.length === 6) {
            setCode(codeParam);
        }

        async function checkCookies() {
            const cookies = new Cookies();

            // If no token cookie, just ignore
            if (!cookies.get('token')) {
                return;
            }

            // If invalid, delete the token
            const res = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (res.status !== 200) {
                console.error('Could not connect to server: error ' + res.status);
                return;
            }
            if (res.data?.ok !== 1) {
                cookies.remove('token');
                return;
            }

            // If all valid, redirect to judge page
            navigate('/judge');
        }

        checkCookies();
    }, [navigate]);

    // Disable button if length of code is not 6
    useEffect(() => {
        setDisabled(code.length === 0);
    }, [code]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCode(e.target.value);
    };

    const handleEnter = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            login();
        }
    };

    const login = async () => {
        // Prevent spamming of login while waiting for fetch
        if (loginLock) return;
        setLoginLock(true);

        // Remove any whitespace in the code
        const cleanCode = code.replace(/\s/g, '');

        // Make async call to check code
        const res = await postRequest<TokenResponse>('/judge/login', '', { code: cleanCode });

        // Invalid code
        if (res.status === 400) {
            setError(true);
            setLoginLock(false);
            return;
        }

        // Internal server error
        if (res.status !== 200) {
            errorAlert(res);

            setError(true);
            setLoginLock(false);
            return;
        }

        // Correct code; save token as cookie
        const token = res.data?.token;
        const cookies = new Cookies();
        cookies.set('token', token, {
            path: '/',
            sameSite: 'strict',
            secure: true,
            maxAge: 60 * 60 * 24,
        });

        // Redirect
        navigate('/judge');

        setLoginLock(false);
    };

    return (
        <>
            <JuryHeader />
            <Container>
                <PasswordInput
                    value={code}
                    label="Enter your judging code"
                    placeholder="XXXXXX"
                    onKeyPress={handleEnter}
                    onChange={handleChange}
                    error={error}
                    setError={setError}
                    errorMessage="Invalid judging code"
                    className="text-3xl md:text-6xl text-center w-4/5"
                />
                <p className="text-2xl text-light text-center mx-4 my-12">
                    If you did not get a code, check your email or contact an organizer.
                </p>
                <Button type="primary" disabled={disabled} onClick={login}>
                    Log In
                </Button>
            </Container>
            <Loading disabled={!loginLock} />
        </>
    );
};

export default JudgeLogin;
