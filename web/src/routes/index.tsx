import { Resource, component$, useBrowserVisibleTask$, useStylesScoped$ } from '@builder.io/qwik';
import type { DocumentHead } from '@builder.io/qwik-city';
import { Link } from '@builder.io/qwik-city';
import { UserTile } from '~/components/UserTile/user-tile';
import { useUsers } from '~/hooks/useUsers';
import type { User } from '~/models/user';
import styles from './index.css?inline';

export default component$(() => {
  useStylesScoped$(styles)

  const { users, reload } = useUsers()

  useBrowserVisibleTask$(() => {

    console.log("we are at client")
    // Get Token from api on the / mockuser route for testing, api has to be startet with USE_MOCK_JWKS = true
    const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IkhhbnMuTWVpZXJAbW9jay5jb20iLCJleHAiOjE2NzY4NDEyOTMsIm5hbWUiOiJIYW5zIE1laWVyIn0.XKgx4E7h27Y96qVIrpcKuMFslaphD_6DqnPBN96KXsQ"

    const eventSource = new EventSource('http://localhost:3333/api/v1/events', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    } as any);

    eventSource.addEventListener("new_question", ({ data }) => {
      console.log("New Question", JSON.parse(data))
    })
    eventSource.addEventListener("update_question", ({ data }) => {
      console.log("Updated Question", JSON.parse(data))
    })
    eventSource.addEventListener("delete_question", ({ data }) => {
      console.log("Deleted Question", JSON.parse(data))
    })
    eventSource.addEventListener("upvote_question", ({ data }) => {
      console.log("Upvoted Question", JSON.parse(data))
    })
    eventSource.addEventListener("answer_question", ({ data }) => {
      console.log("Answered Question", JSON.parse(data))
    })
    eventSource.addEventListener("start_session", () => {
      console.log("Started Session")
    })
    eventSource.addEventListener("stop_session", () => {
      console.log("Stopped Session")
    })

    return () => {
      eventSource.close()
    }
  })

  return (
    <div>
      <button onclick$={reload} type='button'>Reload</button>
      <Resource
        value={users}
        onPending={() => <div>Loading...</div>}
        onRejected={() => <div>Failed to person data</div>}
        onResolved={(users: User[]) => {
          return (
            <div class="user-grid">
              {
                users.map(u =>
                  <Link href={`/user/${u.id}`} key={u.id} class="text-green-600 visited:text-green-600">
                    <UserTile key={u.id} id={u.id} avatar={u.avatar} email={u.email} first_name={u.first_name} last_name={u.last_name} />
                  </Link>
                )
              }
            </div>
          );
        }}
      />
    </div>
  );
});

export const head: DocumentHead = {
  title: 'Welcome to Qwik',
  meta: [
    {
      name: 'description',
      content: 'Qwik site description',
    },
  ],
};

