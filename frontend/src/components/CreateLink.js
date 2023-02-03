import React, { useState } from "react";
import { useMutation, gql } from "@apollo/client";
import { useNavigate } from "react-router-dom";
import { FEED_QUERY } from "./LinkList";
import { AUTH_TOKEN, LINKS_PER_PAGE } from "../constants";

const CREATE_LINK_MUTATION = gql`
  mutation PostMutation($input: NewLink!) {
    post(input: $input) {
      id
      createdAt
      url
      description
    }
  }
`;

const CreateLink = () => {
  const [formState, setFormState] = useState({
    description: "",
    url: "",
  });
  const navigate = useNavigate();
  const [createLink] = useMutation(CREATE_LINK_MUTATION, {
    variables: {
      input: {
        description: formState.description,
        url: formState.url,
      },
    },
    onCompleted: () => navigate("/"),
    update: (cache, { data: { post } }) => {
      const limit = LINKS_PER_PAGE;
      const offset = 0;
      const data = cache.readQuery({ query: FEED_QUERY });
      cache.writeQuery({
        query: FEED_QUERY,
        data: {
          feed: {
            links: [post, ...data.feed.links],
          },
        },
        variables: {
          limit,
          offset,
        }
      });
    },
  });
  return (
    <div>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          createLink();
        }}
      >
        <div className="flex flex-column mt3">
          <input
            className="mb2"
            value={formState.description}
            onChange={(e) =>
              setFormState({
                ...formState,
                description: e.target.value,
              })
            }
            type="text"
            placeholder="A description for the link"
          />
          <input
            className="mb2"
            value={formState.url}
            onChange={(e) =>
              setFormState({
                ...formState,
                url: e.target.value,
              })
            }
            type="text"
            placeholder="The URL for the link"
          />
        </div>
        <button type="submit">Submit</button>
      </form>
    </div>
  );
};

export default CreateLink;
