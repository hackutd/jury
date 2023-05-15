import { useState } from 'react';

const NewJudgeForm = () => {
    const [name, setName] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [notes, setNotes] = useState<string>('');

    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

    const AddJudge = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setIsSubmitting(true);
        try {
            await fetch(`${process.env.REACT_APP_JURY_URL}/judge/new`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, notes }),
                credentials: 'include',
            });
        } catch (err) {
            console.log(err);
            setIsSubmitting(false);
        }

        alert('Judge Added');
        setName('');
        setEmail('');
        setNotes('');
        setIsSubmitting(false);
    };

    return (
        <>
            <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
                <div className="flex flex-col items-start h-full">
                    <h1 className="text-3xl">Add New Judges</h1>
                    <form className="flex flex-col w-full space-y-4" onSubmit={AddJudge}>
                        <div className="flex flex-row w-full mt-4 space-x-6">
                            <input
                                className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                                placeholder="Name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                            <input
                                className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                                placeholder="Email Address"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                        <textarea
                            className="w-full h-36 px-4 py-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                            placeholder="Additional Notes"
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                        />
                        <button className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full">
                            {isSubmitting ? 'Adding...' : 'Add'}
                        </button>
                    </form>
                </div>
            </div>
        </>
    );
};

export default NewJudgeForm;
