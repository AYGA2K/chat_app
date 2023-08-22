<script lang="ts">
	import { goto } from '$app/navigation';

	let name = '';
	let email = '';
	let password = '';
	let confirmPassword = '';

	const signupHandler = async () => {
		if (password !== confirmPassword) {
			alert('passwords do not match');
		}

		const newUser = {
			name,
			email,
			password
		};

		try {
			const response = await fetch('http://localhost:8080/user/signup', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(newUser)
			});

			if (response.ok) {
				goto('/login');
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
			<h1 class="text-5xl font-bold">Sign Up now!</h1>
			<p class="py-6">Join the chat app.</p>
		</div>
		<div class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
			<div class="card-body">
				<div class="form-control">
					<label class="label" for="name">
						<span class="label-text">name</span>
					</label>
					<input bind:value={name} placeholder="name" class="input input-bordered" />
				</div>
				<div class="form-control">
					<label class="label" for="email">
						<span class="label-text">Email</span>
					</label>
					<input type="email" bind:value={email} placeholder="email" class="input input-bordered" />
				</div>
				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						placeholder="password"
						class="input input-bordered"
						bind:value={password}
					/>
				</div>
				<div class="form-control">
					<label class="label" for="confirm-password">
						<span class="label-text">Confirm Password</span>
					</label>
					<input
						type="password"
						bind:value={confirmPassword}
						placeholder="confirm password"
						class="input input-bordered"
					/>
				</div>
				<div class="form-control flex gap-2">
					<button class="btn btn-primary" on:click={signupHandler}>Sign Up</button>
					<a href="/login" class="btn btn-primary">Login</a>
				</div>
			</div>
		</div>
	</div>
</div>
