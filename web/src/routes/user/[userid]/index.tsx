import { component$ } from '@builder.io/qwik';
import { StaticGenerateHandler, useLocation } from '@builder.io/qwik-city';

export default component$(() => {
    const location = useLocation();

    return (
        <div>
            <h1>SKU</h1>
            <p>Pathname: {location.pathname}</p>
            <p>User Id: {location.params.userid}</p>
        </div>
    );
});

export const onStaticGenerate: StaticGenerateHandler = () => {
    const ids = ["1", "2", "3", "4", "5"];

    return {
        params: ids.map((userid) => {
            return { userid };
        }),
    };
};