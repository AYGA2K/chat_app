import { writable } from 'svelte/store';
interface UserData {
	id: any;
	name: any;
}

export const user = writable<UserData | null>(null);
