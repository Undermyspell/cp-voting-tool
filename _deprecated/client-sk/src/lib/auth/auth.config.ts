export const msalConfig = {
	auth: {
		clientId: '7045b40f-1c5e-40da-ab0b-eef6f11d26cb',
		authority: 'https://login.microsoftonline.com/ff930651-2670-491e-9a70-7847e7fbf8b7',
		redirectUri: location.origin
	},
	cache: {
		cacheLocation: 'sessionStorage', // This configures where your cache will be stored
		storeAuthStateInCookie: false // Set this to "true" if you are having issues on IE11 or Edge
	}
};

export const loginRequest = {
	scopes: ['User.Read']
};
