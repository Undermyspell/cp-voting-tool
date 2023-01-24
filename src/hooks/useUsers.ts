import { useResource$, useSignal, $ } from "@builder.io/qwik";
import type { User } from "~/models/user";

export function useUsers() {
	const reloadTrigger = useSignal(true);
	const reload = $(() => (reloadTrigger.value = !reloadTrigger.value));

	const users = useResource$<User[]>(async ({ track }) => {
		track(() => reloadTrigger.value);
		const res = await fetch("https://reqres.in/api/users?page=1");
		const { data }: { data: User[] } = await res.json();
		return data;
	});

	return { users, reload };
}
