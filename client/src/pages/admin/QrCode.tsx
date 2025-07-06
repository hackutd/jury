import qr from 'qr.js';
import { createRef, useEffect, useState } from 'react';
import JuryHeader from '../../components/JuryHeader';
import { useSearchParams } from 'react-router-dom';
import { getRequest, postRequest } from '../../api';
import { errorAlert } from '../../util';

const QrCode = () => {
    const [searchParams, _] = useSearchParams({ track: '' });
    const [code, setCode] = useState('');
    const [track, setTrack] = useState('');
    const canvasRef = createRef<HTMLCanvasElement>();

    useEffect(() => {
        async function fetchData() {
            const tr = searchParams.get('track') ?? '';

            // Get the code
            let cd = '';
            if (tr === '') {
                const res = await getRequest<Code>('/admin/qr', 'admin');
                if (res.status !== 200) {
                    errorAlert(res);
                    return;
                }

                cd = res.data?.qr_code as string;
            } else {
                const res = await getRequest<Code>(`/admin/qr/${tr}`, 'admin');
                if (res.status !== 200) {
                    errorAlert(res);
                    return;
                }

                cd = res.data?.qr_code as string;
            }
            console.log(cd)

            // If code is empty, generate new code
            if (cd === '') {
                generateCode();
            } else {
                setCode(cd);
            }

            setTrack(tr === '' ? 'General' : tr);
        }

        fetchData();
    });

    const generateCode = async () => {
        const track = searchParams.get('track') ?? '';

        // Generate the code
        let cd = '';
        if (track === '') {
            const res = await postRequest<Code>('/admin/qr', 'admin', {});
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }
            cd = res.data?.qr_code as string;
        } else {
            const res = await postRequest<Code>(`/admin/qr/${track}`, 'admin', {});
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }
            cd = res.data?.qr_code as string;
        }
        setCode(cd);
    };

    useEffect(() => {
        if (code === '') return;

        // Create QR code
        const origin = window.location.origin.toString();
        const track = searchParams.get('track') ?? '';
        const url = `${origin}/add-self?code=${code}&track=${track.replace(/\s/g, '%20')}`;

        const qrCode = qr(url);

        // Draw on the canvas
        const ctx = canvasRef.current?.getContext('2d');
        if (ctx) {
            ctx.clearRect(0, 0, 500, 500);
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, 500, 500);

            const size = 500;
            const cellSize = size / qrCode.modules.length;
            qrCode.modules.forEach((row, r) => {
                row.forEach((cell, c) => {
                    ctx.fillStyle = cell ? '#000000' : '#ffffff';
                    ctx.fillRect(c * cellSize, r * cellSize, cellSize, cellSize);
                });
            });
        }
    }, [code]);

    return (
        <>
            <JuryHeader withBack withLogout isAdmin backLocation="/admin/add-judges" />
            <div className="flex items-center flex-col mt-4">
                <h1 className="text-3xl mb-4 text-light font-bold">{track} Judging QR Code</h1>
                <canvas width={500} height={500} ref={canvasRef} />
            </div>
        </>
    );
};

export default QrCode;
