const { S3Client, PutObjectCommand } = require('@aws-sdk/client-s3');

const s3Client = new S3Client({ region: 'ap-southeast-1' });

const S3 = {
    async write(data, fileName, bucket, hash) {
        const params = {
            Bucket: bucket,
            Body: data,
            Key: `${hash}/${fileName}`,
            ContentType: 'application/zip' // Ensure correct content type
        };

        try {
            const command = new PutObjectCommand(params);
            const newData = await s3Client.send(command);

            if (!newData) {
                throw new Error('There was an error writing the file');
            }

            return newData;
        } catch (error) {
            console.error(error);
            throw new Error(`Error writing file: ${error.message}`);
        }
    },
};

module.exports = S3;
