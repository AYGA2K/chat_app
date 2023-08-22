<script lang="ts">
	import { goto } from '$app/navigation';
	import { user } from './../../store';
	let email = '';
	let password = '';

	const loginHandler = async () => {
		const newUser = {
			email,
			password
		};

		try {
			const response = await fetch('http://localhost:8080/user/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(newUser)
			});

			if (response.ok) {
				var userdata = await response.json();
				var User = {
					id: userdata.userID,
					name: userdata.userName
				};
				user.set(User);
				console.log(user);
				goto('/');
			} else {
				console.log(response);
			}
		} catch (error) {
			console.log(error);
		}
	};
</script>

<div class="hero min-h-screen bg-base-200">
	<div class="hero-content flex-col md:w-2/4 lg:flex-row-reverse">
		<div class="text-center lg:text-left">
			<h1 class="text-5xl font-bold">Login now!</h1>
			<p class="py-6">Welcome to the chat app.</p>
		</div>
		<div class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
			<div class="card-body">
				<div class="form-control">
					<label class="label" for="email">
						<span class="label-text">Email</span>
					</label>
					<input type="text" bind:value={email} placeholder="email" class="input input-bordered" />
				</div>
				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						bind:value={password}
						placeholder="password"
						class="input input-bordered"
					/>
				</div>
				<div class="form-control flex gap-2">
					<button on:click={loginHandler} class="btn btn-primary">Login</button>
					<a href="signup" class="btn btn-primary">SignUp</a>
				</div>
			</div>
		</div>
	</div>
</div>
