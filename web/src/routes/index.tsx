import { Resource, component$, useClientEffect$, useStylesScoped$ } from '@builder.io/qwik';
import type { DocumentHead } from '@builder.io/qwik-city';
import { Link } from '@builder.io/qwik-city';
import { UserTile } from '~/components/UserTile/user-tile';
import { useUsers } from '~/hooks/useUsers';
import type { User } from '~/models/user';
import styles from './index.css?inline';

export default component$(() => {
  useStylesScoped$(styles)

  const { users, reload } = useUsers()

  useClientEffect$(() => {

    // Get Token from api on the /mockuser route for testing, api has to be startet with USE_MOCK_JWKS=true
    const token = "-- USER TOKEN --"

    const eventSource = new EventSource('http://localhost:3333/api/v1/events', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    } as any);

    eventSource.addEventListener("new_question", ({ data }) => {
      console.log("New Question", JSON.parse(data))
    })
    eventSource.addEventListener("upvote_question", ({ data }) => {
      console.log("Upvoted Question", JSON.parse(data))
    })
    eventSource.addEventListener("answer_question", ({ data }) => {
      console.log("Answered Question", JSON.parse(data))
    })
    eventSource.addEventListener("reset_session", () => {
      console.log("Reset Session")
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

