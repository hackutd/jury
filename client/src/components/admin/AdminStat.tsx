const AdminStat = (props: { name: string; value: string | number }) => {
    return <div className="text-center">
        <div className="text-6xl">{props.value}</div>
        <div className="text-light">{props.name}</div>
    </div>;
};

export default AdminStat;
