import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, gql } from "@apollo/client";
import { AUTH_TOKEN } from "../constants";

const SIGNUP_MUTATION = gql`
  mutation SignupMutation($input: NewUser!) {
    signup(input: $input) {
      token
    }
  }
`;

const LOGIN_MUTATION = gql`
  mutation LoginMutation($input: Login!) {
    login(input: $input) {
      token
    }
  }
`;

const Login = () => {
  const navigate = useNavigate();
  const [formState, setFormState] = useState({
    login: true,
    email: "",
    password: "",
    username: "",
  });

  const [login] = useMutation(LOGIN_MUTATION, {
    variables: {
      input: {
        username: formState.username,
        password: formState.password,
      },
    },
    onCompleted: ({ login }) => {
      localStorage.setItem(AUTH_TOKEN, login.token);
      navigate("/");
    },
  });
  const [signup] = useMutation(SIGNUP_MUTATION, {
    variables: {
      input: {
        username: formState.username,
        email: formState.email.trim() === "" ? null : formState.email,
        password: formState.password,
      },
    },
    onCompleted: ({ signup }) => {
      localStorage.setItem(AUTH_TOKEN, signup.token);
      navigate("/");
    },
  });

  return (
    <div>
      <h4 className="mv3">{formState.login ? "login" : "Sign up"}</h4>
      <div className="flex flex-column">
        <input
          value={formState.username}
          onChange={(e) =>
            setFormState({
              ...formState,
              username: e.target.value,
            })
          }
          type="text"
          placeholder="Your User name"
        />
        {!formState.login && (
          <input
            value={formState.email}
            onChange={(e) =>
              setFormState({
                ...formState,
                email: e.target.value,
              })
            }
            type="text"
            placeholder="Your email address"
          />
        )}
        <input
          value={formState.password}
          onChange={(e) =>
            setFormState({
              ...formState,
              password: e.target.value,
            })
          }
          type="password"
          placeholder="Choose a safe password"
        />
      </div>
      <div className="flex mt3">
        <button
          className="pointer mr2 button"
          onClick={formState.login ? login : signup}
        >
          {formState.login ? "login" : "create account"}
        </button>
        <button
          className="pointer button"
          onClick={(e) =>
            setFormState({
              ...formState,
              login: !formState.login,
            })
          }
        >
          {formState.login
            ? "need to create an account?"
            : "already have an account?"}
        </button>
      </div>
    </div>
  );
};

export default Login;
