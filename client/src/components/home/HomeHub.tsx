import Button from '../Button';
import Container from '../Container';

const AppHub = () => {
    return (
        <div className="bg-black h-screen w-full">
            <Container>
                <h1 className="text-7xl font-bold text-center mb-4 hover:animate-wiggle cursor-pointer text-white">
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
                    {import.meta.env.VITE_JURY_NAME}
                </h2>
                <Button href="/judge/login" type="primary" large>
                    Judging Portal
                </Button>
                <Button
                    href="/admin/login"
                    type="primary"
                    large
                    className="my-4 bg-primaryLight text-black hover:bg-primaryLight hover:brightness-95 hover:text-black"
                >
                    Admin Portal
                </Button>
                <Button href="/expo" type="outline" className="px-12 py-2 text-xl">
                    Project Expo
                </Button>
            </Container>
        </div>
    );
};

export default AppHub;
