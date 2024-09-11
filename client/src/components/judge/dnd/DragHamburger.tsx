const DragHamburger = () => {
    return (
        <div className="drag-handle p-2 cursor-grab touch-none">
            <svg
                width="28"
                height="25"
                viewBox="0 0 28 25"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className="drag-handle stroke-lightest stroke-[5px]"
            >
                <line y1="2.5" x2="28" y2="2.5" className="drag-handle" />
                <line y1="12.4999" x2="28" y2="12.4999" className="drag-handle" />
                <line y1="22.5" x2="28" y2="22.5" className="drag-handle" />
            </svg>
        </div>
    );
};

export default DragHamburger;
