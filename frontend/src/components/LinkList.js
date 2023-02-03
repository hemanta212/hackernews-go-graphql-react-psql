import React from "react";
import { gql, useQuery } from "@apollo/client";
import { useNavigate, useLocation } from "react-router-dom"
import Link from "./Link";
import { LINKS_PER_PAGE } from "../constants";

export const FEED_QUERY = gql`
  query FeedQuery($limit: Int $offset: Int){
    feed(limit: $limit, offset: $offset) {
      id
      links {
        id
        createdAt
        url
        description
        postedBy {
          id
          username
        }
        votes {
          id
          user{
            id
          }
        }
      }
      count
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
  const navigate = useNavigate();
  const location = useLocation();
  const isNewPage = location.pathname.includes("new");
  const pageIndexParams = location.pathname.split('/');
  const page = parseInt(
    pageIndexParams[pageIndexParams.length - 1]
  )
  const pageIndex = page ? (page - 1) * LINKS_PER_PAGE : 0;
  const { data, loading, error, subscribeToMore } = useQuery(FEED_QUERY, {
    variables: getQueryVariables(isNewPage, page),
  });
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
        ({ id }) => id === newLink.id
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
    <>
      {loading && <p>Loading...</p>}
      {error && <pre>{JSON.stringify(error, null, 2)}</pre>}
      {data && (
        <>
          {getLinksToRender(isNewPage, data).map(
            (link, index) => (
              <Link
                key={link.id}
                link={link}
                index={index + pageIndex}
              />
            )
          )}

          {isNewPage && (
            <div className="flex ml4 mv3 gray">
              <div
                className="pointer mr2"
                onClick={() => {
                  if (page > 1) {
                    navigate(`/new/${page - 1}`)
                  }
                }}
              >prev </div>

              <div
                className="pointer"
                onClick={() => {
                  if (page <= data.feed.count / LINKS_PER_PAGE) {
                    navigate(`/new/${page + 1}`)
                  }
                }}
              >Next </div>
            </div>
          )}
            </>
          )}
        </>
      );
};

const getQueryVariables = (isNewPage, page) => {
  const offset = isNewPage ? (page - 1) * LINKS_PER_PAGE : 0;
      const limit = isNewPage ? LINKS_PER_PAGE : 15;
  //const orderBy = {createdAt: "desc" };
      return {limit, offset};
}

const getLinksToRender = (isNewPage, data) => {
  if (isNewPage){
    return data.feed.links;
  }
  const rankedLinks = data.feed.links.slice();
  rankedLinks.sort(
    (l1, l2) => l2.votes.length - l1.votes.length
  );
  return rankedLinks;
};

      export default LinkList;

