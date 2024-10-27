import { useState, useEffect } from 'react';
import arrow from '../assets/arrow.svg';

function Fab() {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            // show the button when the user has scrolled down more than 200px
            setIsVisible(window.scrollY > 200);
        };

        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    // scroll to the top of the page
    const scrollToTop = () => {
        window.scrollTo({
            top: 0,
            behavior: 'smooth',
        });
    };

    return isVisible ? (
        <button
            onClick={scrollToTop}
            className="fixed bottom-8 right-8 w-12 h-12 border-2 border-primary rounded-md flex items-center justify-center shadow-md bg-white"
        >
            <img src={arrow} alt="arrow up" className="w-5 h-5" />
        </button>
    ) : null;
}

export default Fab;
