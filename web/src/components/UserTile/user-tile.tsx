import { component$ } from '@builder.io/qwik';
import type { User } from '~/models/user';


export const UserTile = component$(({ id, email, first_name, last_name, avatar }: User) => {
    return (
        <div class="flex flex-col gap-2 bg-red-300">
            <img class="w-40 h-40" src={avatar} alt='userimage' />
            <span>{first_name}</span>
            <span>{last_name}</span>
            <span>{email}</span>
            <span>{id}</span>
        </div>
    )
})