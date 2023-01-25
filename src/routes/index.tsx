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
    const eventSource = new EventSource('http://localhost:3333/events');
    eventSource.onmessage = ({ data }) => {
      console.log('New message', JSON.parse(data));
    };
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
                  <UserTile key={u.id} id={u.id} avatar={u.avatar} email={u.email} first_name={u.first_name} last_name={u.last_name} />
                )
              }
            </div>
          );
        }}
      />

      <Link class="mindblow" href="/flower/">
        Blow my mind 🤯
      </Link>
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

