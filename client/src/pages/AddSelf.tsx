import { useNavigate, useSearchParams } from 'react-router-dom';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import { useEffect, useState } from 'react';
import Loading from '../components/Loading';
import TextInput from '../components/TextInput';
import Button from '../components/Button';
import { getRequest, postRequest } from '../api';
import { errorAlert } from '../util';

const AddSelf = () => {
    const [searchParams, _] = useSearchParams({ track: '' });
    const navigate = useNavigate();
    const [loaded, setLoaded] = useState(false);
    const [code, setCode] = useState('');
    const [track, setTrack] = useState('');
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');

    useEffect(() => {
        async function checkCode() {
            // Get track
            const tr = searchParams.get('track') ?? '';

            // Get QR code
            let correctCode;
            if (tr !== '') {
                console.log(tr);
                const res = await getRequest<Code>(`/qr/${tr}`, '');
                if (res.status !== 200) {
                    errorAlert(res);
                    return;
                }

                correctCode = res.data?.qr_code;
            } else {
                const res = await getRequest<Code>('/qr', '');
                if (res.status !== 200) {
                    errorAlert(res);
                    return;
                }
                correctCode = res.data?.qr_code;
            }

            // Check code
            const paramCode = searchParams.get('code');
            if (!paramCode || paramCode !== correctCode) {
                alert('Invalid code, please re-scan the QR code or ask an organizer.');
                return;
            }

            setCode(paramCode);
            setTrack(tr);
            setLoaded(true);
        }

        checkCode();
    }, []);

    const createJudge = async () => {
        setLoaded(false);

        const res = await postRequest<OkResponse>('/qr/add', '', {
            name,
            email,
            track,
            code,
        });
        if (res.status !== 200) {
            errorAlert(res);
            setLoaded(true);
            return;
        }

        navigate('/add-self/done');

        setLoaded(true);
    };

    return (
        <>
            <JuryHeader />
            <Container>
                <h1 className="text-3xl">Add Judge Form</h1>
                <p className="text-light mt-2 mb-8 px-4 text-center">
                    Enter your information below to add yourself to the judging system -- all you
                    need is your name and email. Once you hit submit, you will get an email with
                    your judging code.
                </p>
                <TextInput text={name} setText={setName} label="Name" large />
                <TextInput
                    text={email}
                    setText={setEmail}
                    label="Email"
                    large
                    className="mt-2 mb-8"
                />
                <Button type="primary" onClick={createJudge}>
                    Submit
                </Button>
                <Loading disabled={loaded} />
            </Container>
        </>
    );
};

export default AddSelf;
