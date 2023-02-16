import JuryHeader from '../components/JuryHeader';

const Admin = () => {
    return (
        <>
            <JuryHeader withLogout />
            <div className="panel">
                <div className="stats-row">
                    <div className="stat">
                        <p id="seen" className="stat-title">
                            0
                        </p>
                        <p className="stat-subtitle">Seen</p>
                    </div>
                    <div className="stat">
                        <p id="votes" className="stat-title">
                            0
                        </p>
                        <p className="stat-subtitle">Votes</p>
                    </div>
                    <div className="stat">
                        <p id="time" className="stat-title">
                            0
                        </p>
                        <p className="stat-subtitle">Judging Time</p>
                    </div>
                    <div className="stat">
                        <p id="mu" className="stat-title">
                            0
                        </p>
                        <p className="stat-subtitle">Average Mu</p>
                    </div>
                    <div className="stat">
                        <p id="sigma" className="stat-title">
                            0
                        </p>
                        <p className="stat-subtitle">Average Sigma^2</p>
                    </div>
                </div>
                <div className="control-row">
                    <div className="batch">
                        <p>Batch Ops</p>
                        <div className="batch-icon">Prioritize</div>
                        <div className="batch-icon">Edit</div>
                        <div className="batch-icon">Hide</div>
                        <div className="batch-icon">Delete</div>
                    </div>
                    <div className="tabs">
                        <div className="tab-text"></div>
                        <div className="tab-text"></div>
                        <div className="tab-text"></div>
                        <div className="tab-rect"></div>
                    </div>
                    <div className="add-container">
                        <p className="add-judges">Add Judges</p>
                    </div>
                </div>
                <table>
                    <tr>
                        <th>Name</th>
                        <th>Votes</th>
                        <th>Alpha</th>
                        <th>Beta</th>
                        <th>Updated</th>
                        <th>Actions &lbrack;Delete,Hide,Edit&rbrack;</th>
                    </tr>
                </table>
            </div>
        </>
    );
};

export default Admin;
