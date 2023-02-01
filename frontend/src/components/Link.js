import React from "react";
import { gql, useMutation } from "@apollo/client";
import { AUTH_TOKEN } from "../constants";
import { timeDifferenceForDate } from "../utils";
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
  const [upVote] = useMutation(VOTE_MUTATION, {
    variables: {
      linkID: link.id,
    },
    update: (cache, { data: { vote } }) => {
      const { feed } = cache.readQuery({ query: FEED_QUERY });
      console.log(feed);
      const updatedLinks = feed.links.map((feedLink) => {
        console.log(feedLink.votes.length);
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
            onClick={upVote}
          >
            ▲
          </div>
        )}
      </div>
      <div className="ml1">
        <div>
          {link.description} ({link.url})
        </div>
        {
          <div className="f6 1h-copy gray">
            {link.votes.length} votes | by{" "}
            {link.postedBy ? link.postedBy.username : "Anonymous"}{" "}
            {timeDifferenceForDate(link.createdAt)}
          </div>
        }
      </div>
    </div>
  );
};

export default Link;
