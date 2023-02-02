import React from "react";
import { gql, useQuery } from "@apollo/client";
import Link from "./Link";

export const FEED_QUERY = gql`
  {
    feed {
      links {
        id
        createdAt
        description
        url
        postedBy {
          username
        }
        votes {
          id
        }
      }
    }
  }
`;

export const NEW_LINKS_SUBSCRIPTION = gql`
  subscription {
    newLink {
      id  
      url 
      description
      createdAt
      postedBy {
        id
        username
      }
      votes {
        id
        user {
          id
        }
      }
    }
  }
`;

export const NEW_VOTES_SUBSCRIPTION = gql`
  subscription {
    newVote {
      id
      link {
        id
        url
        description
        createdAt
        postedBy {
          id
          username
        }
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

const LinkList = () => {
  const { data, loading, error, subscribeToMore } = useQuery(FEED_QUERY);
  subscribeToMore({
      document: NEW_VOTES_SUBSCRIPTION,
  });
  subscribeToMore({
      document: NEW_LINKS_SUBSCRIPTION,
      // updateQuery is a function that takes the previous state of the query
      // and the subscription data and returns the new state of the query
      // these args are passed in automatically by Apollo
      updateQuery: (prev, { subscriptionData }) => {
        if (!subscriptionData.data) return prev;
        const newLink = subscriptionData.data.newLink;
        // find is a method on arrays that returns the first item that matches
        // the condition in the callback
        console.log("curr id prev links ids: ", newLink.id, prev.feed.links.map(link => link.id))
        const exists = prev.feed.links.find(
          ({id}) => id === newLink.id
        );
        if (exists) return prev;

        return Object.assign({}, prev, {
          feed: {
            links: [newLink, ...prev.feed.links],
            count: prev.feed.links.length + 1,
            __typename: prev.feed.__typename
        }
      });
      }
  });
  return (
    <div>
      {data && (
        <>
          {data.feed.links.map((link, index) => (
            <Link key={link.id} link={link} index={index} />
          ))}
        </>
      )}
    </div>
  );
};

export default LinkList;
