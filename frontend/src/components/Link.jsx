import React from "react";
import { AUTH_TOKEN, LINKS_PER_PAGE } from "../constants";
import { timeDifferenceForDate } from "../utils";
import { useMutation, gql } from "@apollo/client";
import { FEED_QUERY } from "./LinkList";

const VOTE_MUTATION = gql`
  mutation VoteMutation($linkID: ID!) {
    vote(linkID: $linkID) {
      id
      link {
        id
        votes {
          id
          user {
            id
          }
        }
      }
      user {
        id
      }
    }
  }
`;

const Link = (props) => {
  const { link } = props;
  const authToken = localStorage.getItem(AUTH_TOKEN);

  const limit = LINKS_PER_PAGE;
  const offset = 0;
  const orderBy = { createdAt: "desc" };

  const [vote] = useMutation(VOTE_MUTATION, {
    variables: {
      linkID: link.id,
    },
    update: (cache, { data: { vote } }) => {
      const { feed } = cache.readQuery({
        query: FEED_QUERY,
        variables: {
          limit: limit,
          offset: offset,
          orderBy,
        },
      });

      const updatedLinks = feed.links.map((feedLink) => {
        if (feedLink.id === link.id) {
          return {
            ...feedLink,
            votes: [...feedLink.votes, vote],
          };
        }
        return feedLink;
      });

      cache.writeQuery({
        query: FEED_QUERY,
        data: {
          feed: {
            links: updatedLinks,
          },
        },
        variables: {
          limit: limit,
          offset: offset,
          orderBy,
        },
      });
    },
  });

  return (
    <div className="flex mt2 items-start">
      <div className="flex items-center">
        <span className="gray">{props.index + 1}.</span>
        {authToken && (
          <div
            className="ml1 gray f11"
            style={{ cursor: "pointer" }}
            onClick={() => {
              vote();
              props.refetch();
            }}
          >
            ▲
          </div>
        )}
      </div>
      <div className="ml1">
        <div>
          {link.description} (<a href={link.url}>{link.url}</a>)
        </div>
        {
          <div className="f6 lh-copy gray">
            {link.votes.length} votes | by{" "}
            {link.postedBy ? link.postedBy.username : "Unknown"}{" "}
            {timeDifferenceForDate(link.createdAt)}
          </div>
        }
      </div>
    </div>
  );
};

export default Link;
