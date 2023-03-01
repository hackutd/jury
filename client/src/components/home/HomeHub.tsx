import Button from '../Button';
import Container from '../Container';

const AppHub = () => {
    return (
        <div className="bg-black h-full">
            <Container>
                <h1 className="text-7xl font-bold text-center mb-8 hover:animate-wiggle cursor-pointer text-white">
                    <a
                        href="https://github.com/acmutd/jury"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Judge
                        <span className="bg-[#FF9839] rounded text-black ml-2 px-2">Hub</span>
                    </a>
                </h1>
                <h2 className="text-primary text-3xl text-center font-bold mb-24">
                    {process.env.REACT_APP_JURY_NAME}
                </h2>
                <Button href="/judge/login" type="primary">
                    Judging Portal
                </Button>
                <Button href="/admin/login" type="text">
                    Admin Portal
                </Button>
            </Container>
        </div>
    );
};

export default AppHub;
