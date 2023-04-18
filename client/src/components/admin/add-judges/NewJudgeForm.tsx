const NewJudgeForm = () => {
    return (
        <>
            <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
                <div className="flex flex-col items-start h-full">
                    <h1 className="text-3xl">Add New Judges</h1>
                    <form className="flex flex-col w-full space-y-4">
                        <div className="flex flex-row w-full mt-4 space-x-6">
                            <input
                                className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                                placeholder="Name"
                            />
                            <input
                                className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                                placeholder="Email Address"
                            />
                        </div>
                        <textarea
                            className="w-full h-36 px-4 py-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                            placeholder="Additional Notes"
                        />
                        <button className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full">
                            Add
                        </button>
                    </form>
                </div>
            </div>
        </>
    );
};

export default NewJudgeForm;
