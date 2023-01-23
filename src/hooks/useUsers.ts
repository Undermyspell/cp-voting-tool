import { useResource$, useSignal } from "@builder.io/qwik";
import { User } from "~/models/user";

export const loadUsers = async (): Promise<User[]> => {
	const res = await fetch("https://reqres.in/api/users?page=1");
	const { data }: { data: User[] } = await res.json();
	return data;
};

export function useUsers() {
	const users = useResource$<User[]>(async () => {
		return await loadUsers();
	});

	return { users };
}
