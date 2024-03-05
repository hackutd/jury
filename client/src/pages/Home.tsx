import Button from '../components/Button';
import Container from '../components/Container';

const App = () => {
    return (
        <Container>
            <h1 className="text-7xl font-bold text-center mb-4 hover:animate-wiggle cursor-pointer">
                <a href="https://github.com/acmutd/jury" target="_blank" rel="noopener noreferrer">
                    Jury
                </a>
            </h1>
            <h2 className="text-primary text-3xl text-center font-bold mb-24">
                {import.meta.env.VITE_JURY_NAME}
            </h2>
            <Button href="/judge/login" type="primary">
                Judging Portal
            </Button>
            <Button href="/expo" type="outline" className="py-3 mt-4 mb-2">
                Project Expo
            </Button>
            <Button href="/admin/login" type="text">
                Admin Portal
            </Button>
        </Container>
    );
};

export default App;
