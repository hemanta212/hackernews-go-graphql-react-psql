import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { AUTH_TOKEN } from "../constants";

const Header = () => {
  const navigate = useNavigate();
  const authToken = localStorage.getItem(AUTH_TOKEN);
  return (
    <div className="flex pa1 justfiy-between nowrap orange">
      <div className="flex flex-fixed black">
        <Link to="/" className="no-underline black">
          <div className="fw7 mr1">Hacker News </div>
        </Link>
        <Link to="/top" className="ml1 no-underline black">
          top
        </Link>
        <div className="ml1">|</div>
        <Link to="/" className="ml1 no-underline black">
          new
        </Link>
        <div className="ml1">|</div>
        <Link to="/search" className="ml1 no-underline black">
          search
        </Link>
        {authToken && (
          <div className="flex">
            <div className="ml1">|</div>
            <Link to="/create" className="ml1 no-underline black">
              submit
            </Link>
          </div>
        )}
      </div>
      <div className="flex tr-1 flex-fixed">
        {authToken ? (
          <div
            className="ml2 pointer black"
            onClick={() => {
              localStorage.removeItem(AUTH_TOKEN);
              navigate("/");
            }}
          >
            | logout
          </div>
        ) : (
          <Link to="/login" className="ml2 no-underline black">
            login
          </Link>
        )}
      </div>
    </div>
  );
};

export default Header;
