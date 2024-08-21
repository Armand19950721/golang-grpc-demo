

(async () => {
    // import
    const utils = require('./utils')
    const { Client } = require('pg')

    // connect psql
    const client = new Client({
        user: "spe3d", // string | undefined;
        database: "meta_commerce", // string | undefined;
        password: "", // string | (() => string | Promise<string>) | undefined;
        port: "5432", // number | undefined;
        host: "", // string | undefined;
    })
    await client.connect()

    // get data
    const res = await client.query('select * from ar_content')
    const datas = res.rows

    // loop rows and update each one
    for (let i = 0; i < datas.length; i++) {
        const item = datas[i];
        let need_update = false;

        // try to get nosql col
        if (!!item.viewer_setting) {
            let obj = JSON.parse(item.viewer_setting);
            console.log({ obj })

            // replace col
            utils.replaceColumn(obj, 'viwer_right_button', 'viewer_right_button')
            utils.replaceColumn(obj, 'viwer_left_button', 'viewer_left_button')
            utils.replaceColumn(obj, 'button_color', 'camera_button_color')
            console.log({ new: obj })

            // replace obj
            item.viewer_setting = JSON.stringify(obj)
            need_update = true
        }

        // update if needed
        if (need_update) {
            console.log({ item })
            const update = await client.query(`UPDATE ar_content SET viewer_setting = '${item.viewer_setting}' WHERE id = '${item.id}';`)
            console.log({ update })
        }
    }

    return
})();