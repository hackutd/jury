const UploadCSVForm = () => {
    return (
        <>
            <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
                <div className="flex flex-col items-start h-full">
                    <h1 className="text-3xl">Upload CSV</h1>
                    <p className="text-lg text-light">
                        CSV should be formatted in the order of name, email, and notes separated by
                        commas and then data on newlines
                    </p>
                    <form className="flex flex-col w-full space-y-4 mt-4">
                        <div className="flex items-center justify-center w-full">
                            <label
                                htmlFor="dropzone-file"
                                className="flex flex-col items-center justify-center w-full h-64 border-2 border-primary border-dashed rounded-sm cursor-pointer bg-primary/20"
                            >
                                <div className="flex flex-col items-center justify-center pt-5 pb-6">
                                    <p className="text-2xl text-center">
                                        Drag & Drop files here to upload, <br />
                                        or <span className="text-primary"> Browse Files</span>
                                    </p>
                                </div>
                                <input id="dropzone-file" type="file" className="hidden" />
                            </label>
                        </div>
                        <div className="flex flex-row w-full mt-4 space-x-6">
                            <button className="w-full h-11 px-4 text-2xl text-light border-2 border-lightest rounded-full">
                                Cancel
                            </button>
                            <button className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full">
                                Upload
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </>
    );
};

export default UploadCSVForm;
